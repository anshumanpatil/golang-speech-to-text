# pip freeze > requirements.txt
# pip uninstall -r requirements.txt -y
# pip install -r requirements.txt
from pocketsphinx import LiveSpeech, get_model_path
# from kafka import KafkaProducer
# import translators.server as tss
# import translators as ts
import os

MODELDIR = os.path.join(os.path.dirname(__file__), "cmusphinx-en-in-8khz-5.2") # General english model
# MODELDIR = os.path.join(os.path.dirname(__file__), "7088") # Custom words model

print("MODELDIR ", MODELDIR)
# bootstrap_servers = ['localhost:9092']
# topicName = 'speech'
# producer = KafkaProducer(bootstrap_servers = bootstrap_servers)
# producer = KafkaProducer()

# http://www.speech.cs.cmu.edu/tools/lmtool-new.html
# ffmpeg is must

HMM_PATH = MODELDIR + "/en_in.cd_cont_5000" 
LM_PATH = MODELDIR + "/en-us.lm.bin"
DIC_PATH = MODELDIR + "/en_in.dic"

def process_question(phrase):
#     # producer.send(topicName, phrase)
#     if language is not None:
#         translation = ts.translate_text(str(phrase), if_ignore_empty_query=True, if_ignore_limit_of_length=True, to_language=language)
#         print(phrase, " - ", translation)

    print("phrase -", phrase, "-")

speech = LiveSpeech(
    sampling_rate=16000,
    hmm=HMM_PATH,
    lm=LM_PATH,
    dic=DIC_PATH
)

for phrase in speech: 
    # print(phrase.segments(detailed=True))
    process_question(phrase)



