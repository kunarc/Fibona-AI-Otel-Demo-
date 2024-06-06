import grpc
from concurrent import futures

import sentence_transformers
from langchain.chains.retrieval_qa.base import RetrievalQA
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain_community.vectorstores import FAISS
from langchain_text_splitters import RecursiveCharacterTextSplitter
from openai import BaseModel

import proto.chat_pb2 as chat_pb2
import proto.chat_pb2_grpc as chat_pb2_grpc
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.grpc import server_interceptor
from opentelemetry.sdk.resources import Resource
from opentelemetry.semconv.resource import ResourceAttributes
from llm_class import MyQwen
from opentelemetry.propagate import extract
model_name = r"./emb/maidalun/bce-embedding-base_v1"
model_kwargs = {'device': 'cpu'}
encode_kwargs = {'normalize_embeddings': False}
embeddings = HuggingFaceEmbeddings(
    model_name=model_name,
    model_kwargs=model_kwargs,
    encode_kwargs=encode_kwargs
)
embeddings.client = sentence_transformers.SentenceTransformer(embeddings.model_name, device='cpu')
# 初始化加载器
text_splitter = RecursiveCharacterTextSplitter(chunk_size=256, chunk_overlap=128)


class Query(BaseModel):
    text: str

class ChatService(chat_pb2_grpc.SendChatServicer):
    def SendChat(self, request, context):
        tracer = trace.get_tracer("sendChat-tracer")
        # 从 gRPC 请求的 context 中提取上下文信息
        metadata = dict(context.invocation_metadata())
        parent_context = extract(metadata)
        # parent_context = extract(lambda carrier, key: context.invocation_metadata().get(key),
        #                          context.invocation_metadata())

        # 创建一个新的 Span，继承 gRPC 请求的上下文
        with tracer.start_as_current_span("SendChat", context=parent_context) as span:
            # 设置 Span 的属性
            span.set_attribute("request-llm-message", request.message)
            llm = MyQwen()
            span.add_event("start-select")
            qa = select_docs(llm, span)
            span.add_event("chatWithLLM")
            res = qa.invoke(request.message)['result']
            # 设置响应结果作为 Span 的属性
            span.set_attribute("response-llm-reply", res)
            span.add_event("server-complete")
            return chat_pb2.ChatResponse(response=res)
        # llm = MyQwen()
        # qa = select_docs(llm)
        # res = qa.invoke(request.message)['result']
        # return chat_pb2.ChatResponse(response=res)

def build_knowledge_base(text):
    knowledge = text
    split_docs = text_splitter.split_text(knowledge)
    db = FAISS.load_local('./faiss', embeddings=embeddings, allow_dangerous_deserialization=True)
    db.add_texts(split_docs)
    db.save_local("./faiss")

# 选取文档
def select_docs(llm, parent_span):
    tracer = trace.get_tracer("select-tracer")
    with tracer.start_as_current_span("SelectDocs", context=trace.set_span_in_context(parent_span)) as span:
        span.add_event("Loading FAISS index")
        db = FAISS.load_local('./faiss', embeddings=embeddings, allow_dangerous_deserialization=True)
        span.add_event("FAISS index loaded")
        retriever = db.as_retriever()
        span.add_event("Retriever created")
        qa = RetrievalQA.from_chain_type(llm=llm, chain_type="stuff", retriever=retriever)
        span.add_event("QA chain created")
        return qa
def serve():
    resource = Resource.create({ResourceAttributes.SERVICE_NAME: "py-server-otel", ResourceAttributes.SERVICE_VERSION: "v0.1.0"})
    # 配置 OpenTelemetry
    tracer_provider = TracerProvider(resource = resource)
    trace.set_tracer_provider(tracer_provider)

    # 配置 OTLP 导出器
    otlp_exporter = OTLPSpanExporter(endpoint="localhost:4317", insecure=True)
    span_processor = BatchSpanProcessor(otlp_exporter)
    tracer_provider.add_span_processor(span_processor)


    # 配置 gRPC 服务器
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), interceptors=(server_interceptor(),))
    chat_pb2_grpc.add_SendChatServicer_to_server(ChatService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()