from pocketsphinx import LiveSpeech
# from pocketsphinx import Decoder
import os
import pyaudio
# http://www.speech.cs.cmu.edu/tools/lmtool-new.html

# MODELDIR = "/Users/apple/projects/ai/pocketsphinx/model" 
# DATADIR = "/Users/apple/projects/ai/pocketsphinx/data"
def process_question(phrase):
    # if str(phrase) == "hello":
        # print(phrase, " User, How are you?")

    print(phrase)

speech = LiveSpeech(
    sampling_rate=16000,  # optional
    hmm='/Users/apple/Downloads/cmusphinx-en-in-8khz-5.2/en_in.cd_cont_5000',
    lm='/Users/apple/Downloads/cmusphinx-en-in-8khz-5.2/en-us.lm.bin',
    dic='/Users/apple/Downloads/cmusphinx-en-in-8khz-5.2/en_in.dic'
)
# speech = LiveSpeech(
#     sampling_rate=16000,  # optional
#     hmm='/Users/apple/projects/ai/pocketsphinx/model/en-us/en-us',
#     lm='/Users/apple/Downloads/TAR6901/6901.lm',
#     dic='/Users/apple/Downloads/TAR6901/6901.dic'
# )
for phrase in speech: 
    # print(phrase.segments(detailed=True))
    process_question(phrase)




