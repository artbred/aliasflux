import json
import vertexai
import google.cloud.aiplatform as aiplatform

from google.oauth2 import service_account
from langchain.chat_models import ChatVertexAI
from langchain.prompts import ChatPromptTemplate
from langchain.schema.messages import HumanMessage, AIMessage
from langchain.chat_models

with open("service_account.json", encoding="utf-8") as f:
    service_account_info = json.load(f)
    project_id = service_account_info["project_id"]

    cred = service_account.Credentials.from_service_account_info(service_account_info)
    vertexai.init(project=project_id, location="us-central1")
    aiplatform.init(credentials=cred)


chat = ChatVertexAI(
    model_name='chat-bison@001',
    temperature=0.4,
    max_output_tokens=512,
    top_p=0.95,
    top_k=40,
    request_parallelism=10
)


def send_chat_message(template: str, previous_messages: list = None, **kwargs):
    all_messages = []
    
    if previous_messages is not None:
        for message in previous_messages:
            sender = message["sender"]
            del message["sender"]

            if sender == "ai":
                all_messages.append(AIMessage(content=str(message)))
            elif sender == "user":
                all_messages.append(HumanMessage(content=str(message)))

    # I did not manage it to work in time
    all_messages.reverse()

    prompt_template = ChatPromptTemplate.from_template(template)

    messages = prompt_template.format_messages(**kwargs)
    all_messages.extend(messages)

    return chat(messages).content
