import os
import json
from flask import Flask, Response, abort, request
from ibm_watson import NaturalLanguageUnderstandingV1
from ibm_cloud_sdk_core.authenticators import IAMAuthenticator

# Set up logging
import logging
from logging import StreamHandler

# Define the base logger
logging.getLogger("analyzer").setLevel(logging.DEBUG)
log = logging.getLogger("analyzer")
stream_handler = StreamHandler()
stream_formatter = logging.Formatter('[%(asctime)s] [%(thread)d] [%(module)s : %(lineno)d] [%(levelname)s] %(message)s')
stream_handler.setFormatter(stream_formatter)
log.addHandler(stream_handler)

# Flask configuration
app = Flask(__name__, static_url_path='')
app.config['PROPAGATE_EXCEPTIONS'] = True

# Define global variables for Watson API
global tone_analyzer_ep
global api_key

def analyze_tone(input_text):
    """
    Function to analyze tone using IBM Watson Natural Language Understanding API.

    :param input_text: Text to be analyzed for tone.
    :return: List containing tone information.
    """
    api_key = os.getenv('TONE_ANALYZER_API_KEY')
    api_url = os.getenv('TONE_ANALYZER_SERVICE_API')
    authenticator = IAMAuthenticator(api_key)
    
    # Create Natural Language Understanding instance
    natural_language_understanding = NaturalLanguageUnderstandingV1(
        version='2021-08-01',
        authenticator=authenticator
    )
    natural_language_understanding.set_service_url(api_url)
    
    # Analyze the input text for tone
    response = natural_language_understanding.analyze(
        text=input_text,
        features={
            "classifications": {
                "model": "tone-classifications-en-v1",
            }
        },
    )
    
    # Extract and return tone information
    tone = response.result.get("classifications")[0].get("class_name")
    return [{ "tone_name" : tone }]

# API endpoint to analyze tone
@app.route('/tone', methods=['POST'])
def tone():
    """
    API endpoint to analyze tone.
    Expects POST request with JSON body containing 'input_text'.

    :return: JSON response with tone analysis.
    """
    log.info("Serving /tone")
    
    # Check if request contains valid JSON body
    if not request.json or not 'input_text' in request.json:
        log.error("bad body: '%s'", request.data)
        abort(400)
    
    # Extract input_text from request
    input_text = request.json['input_text']
    log.info("input text is '%s'", input_text)
    
    # Analyze tone and return JSON response
    tone_doc = analyze_tone(input_text)
    return (json.dumps(tone_doc), 200)

# Home route
@app.route('/', methods=['GET'])
def home():
    """
    Home route for testing if the application is running.

    :return: Simple message indicating application is running.
    """
    log.info("home")
    return ("This works")

# Main function to run the Flask application
if __name__ == '__main__':
    PORT = '5000'
    app.run(host='0.0.0.0', port=int(PORT))
