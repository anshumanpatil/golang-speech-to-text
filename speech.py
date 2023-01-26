#!/usr/bin/env python3

import time
from kafka import KafkaProducer
import speech_recognition as sr


bootstrap_servers = ['localhost:9092']
topicName = 'speech'
producer = KafkaProducer(bootstrap_servers = bootstrap_servers)
producer = KafkaProducer()

def callback(recognizer, audio):
    try:
        # print("Google Speech Recognition thinks you said " + recognizer.recognize_google(audio))
        res = bytes(recognizer.recognize_google(audio), 'utf-8')
        producer.send(topicName, res)
        # metadata = ack.get()
        # print("metadata.topic "+metadata.topic)
        # print("metadata.partition "+str(metadata.partition))
    except sr.UnknownValueError:
        print("Google Speech Recognition could not understand audio")
    except sr.RequestError as e:
        print("Could not request results from Google Speech Recognition service; {0}".format(e))


r = sr.Recognizer()
m = sr.Microphone()
with m as source:
    r.adjust_for_ambient_noise(source)
    print("Started")

stop_listening = r.listen_in_background(m, callback)
for _ in range(50): time.sleep(0.1)  
while True: time.sleep(0.1)  