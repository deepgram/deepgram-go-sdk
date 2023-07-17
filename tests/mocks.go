package deepgram_test

const MockAPIKey = "m0ckap1k3y0bbc125dac7f40ed3eb0ed232a2ff8"
const MockRequestId = "44617f75-5053-4fb1-a30d-7714eee9d414"
const MockAudioURL = "https://fake.deepgram.test.mockaudio.wav"

var MockBasicPreRecordedResponse = `{
	"metadata": {
	  "transaction_key": "deprecated",
	  "request_id": "700df3f5-1180-4434-a089-f3d71e854c5c",
	  "sha256": "5c3d7d6564176747187095038090a5bb44e703968c3836cfba29bf2c0c7d5dbe",
	  "created": "2023-07-03T19:20:26.829Z",
	  "duration": 21.687187,
	  "channels": 1,
	  "models": [
		"96a295ec-6336-43d5-b1cb-1e48b5e6d9a4"
	  ],
	  "model_info": {
		"96a295ec-6336-43d5-b1cb-1e48b5e6d9a4": {
		  "name": "general",
		  "version": "2023-02-22.3",
		  "arch": "base"
		}
	  }
	},
	"results": {
	  "channels": [
		{
		  "alternatives": [
			{
			  "transcript": "my illness has made very little difference to my scientific work it was likely to have chosen theoretical physics that is mainly your thought for which my physical disability is no handicap i may be mentally disabled as well but if so i am too far long to realize",
			  "confidence": 0.9965741,
			  "words": [
				{
				  "word": "my",
				  "start": 0.439466,
				  "end": 0.6791747,
				  "confidence": 0.9941282
				},
				{
				  "word": "illness",
				  "start": 0.6791747,
				  "end": 1.1585922,
				  "confidence": 0.99941206
				},
				{
				  "word": "has",
				  "start": 1.1585922,
				  "end": 1.3983009,
				  "confidence": 0.995789
				},
				{
				  "word": "made",
				  "start": 1.3983009,
				  "end": 1.7179126,
				  "confidence": 0.9998864
				},
				{
				  "word": "very",
				  "start": 1.7179126,
				  "end": 2.0375242,
				  "confidence": 0.99977845
				},
				{
				  "word": "little",
				  "start": 2.0375242,
				  "end": 2.4370387,
				  "confidence": 0.9998828
				},
				{
				  "word": "difference",
				  "start": 2.4370387,
				  "end": 2.9164562,
				  "confidence": 0.9888885
				},
				{
				  "word": "to",
				  "start": 2.9164562,
				  "end": 3.076262,
				  "confidence": 0.9215757
				},
				{
				  "word": "my",
				  "start": 3.076262,
				  "end": 3.3958735,
				  "confidence": 0.9985771
				},
				{
				  "word": "scientific",
				  "start": 3.3958735,
				  "end": 3.8958735,
				  "confidence": 0.99998593
				},
				{
				  "word": "work",
				  "start": 4.194903,
				  "end": 4.4346113,
				  "confidence": 0.9972281
				},
				{
				  "word": "it",
				  "start": 4.914029,
				  "end": 5.1537375,
				  "confidence": 0.98879015
				},
				{
				  "word": "was",
				  "start": 5.1537375,
				  "end": 5.3934464,
				  "confidence": 0.99965894
				},
				{
				  "word": "likely",
				  "start": 5.3934464,
				  "end": 5.792961,
				  "confidence": 0.8131681
				},
				{
				  "word": "to",
				  "start": 5.792961,
				  "end": 6.0326695,
				  "confidence": 0.9951219
				},
				{
				  "word": "have",
				  "start": 6.0326695,
				  "end": 6.352281,
				  "confidence": 0.91999376
				},
				{
				  "word": "chosen",
				  "start": 6.352281,
				  "end": 6.831699,
				  "confidence": 0.99996305
				},
				{
				  "word": "theoretical",
				  "start": 6.831699,
				  "end": 7.331699,
				  "confidence": 0.9999466
				},
				{
				  "word": "physics",
				  "start": 7.630728,
				  "end": 8.110146,
				  "confidence": 0.9998431
				},
				{
				  "word": "that",
				  "start": 8.762314,
				  "end": 8.921488,
				  "confidence": 0.9996525
				},
				{
				  "word": "is",
				  "start": 8.921488,
				  "end": 9.160248,
				  "confidence": 0.9078474
				},
				{
				  "word": "mainly",
				  "start": 9.160248,
				  "end": 9.660248,
				  "confidence": 0.9965741
				},
				{
				  "word": "your",
				  "start": 9.717356,
				  "end": 10.035703,
				  "confidence": 0.9974491
				},
				{
				  "word": "thought",
				  "start": 10.035703,
				  "end": 10.535703,
				  "confidence": 0.9694835
				},
				{
				  "word": "for",
				  "start": 10.83157,
				  "end": 11.070331,
				  "confidence": 0.99951303
				},
				{
				  "word": "which",
				  "start": 11.070331,
				  "end": 11.388678,
				  "confidence": 0.99987257
				},
				{
				  "word": "my",
				  "start": 11.388678,
				  "end": 11.627438,
				  "confidence": 0.9997918
				},
				{
				  "word": "physical",
				  "start": 11.627438,
				  "end": 12.127438,
				  "confidence": 0.9998179
				},
				{
				  "word": "disability",
				  "start": 12.5824795,
				  "end": 13.0824795,
				  "confidence": 0.9998554
				},
				{
				  "word": "is",
				  "start": 13.378347,
				  "end": 13.617107,
				  "confidence": 0.9034617
				},
				{
				  "word": "no",
				  "start": 13.617107,
				  "end": 14.117107,
				  "confidence": 0.9994584
				},
				{
				  "word": "handicap",
				  "start": 14.412975,
				  "end": 14.651735,
				  "confidence": 0.9570983
				},
				{
				  "word": "i",
				  "start": 15.208843,
				  "end": 15.447603,
				  "confidence": 0.99884593
				},
				{
				  "word": "may",
				  "start": 15.447603,
				  "end": 15.686363,
				  "confidence": 0.9965539
				},
				{
				  "word": "be",
				  "start": 15.686363,
				  "end": 15.925123,
				  "confidence": 0.99631906
				},
				{
				  "word": "mentally",
				  "start": 15.925123,
				  "end": 16.425123,
				  "confidence": 0.97878766
				},
				{
				  "word": "disabled",
				  "start": 16.720991,
				  "end": 17.118925,
				  "confidence": 0.99963963
				},
				{
				  "word": "as",
				  "start": 17.118925,
				  "end": 17.357685,
				  "confidence": 0.99978465
				},
				{
				  "word": "well",
				  "start": 17.357685,
				  "end": 17.857685,
				  "confidence": 0.99977034
				},
				{
				  "word": "but",
				  "start": 18.079582,
				  "end": 18.2375,
				  "confidence": 0.9989341
				},
				{
				  "word": "if",
				  "start": 18.2375,
				  "end": 18.553333,
				  "confidence": 0.99781287
				},
				{
				  "word": "so",
				  "start": 18.553333,
				  "end": 18.790207,
				  "confidence": 0.9955264
				},
				{
				  "word": "i",
				  "start": 19.185,
				  "end": 19.263958,
				  "confidence": 0.5573255
				},
				{
				  "word": "am",
				  "start": 19.263958,
				  "end": 19.500834,
				  "confidence": 0.8580289
				},
				{
				  "word": "too",
				  "start": 19.500834,
				  "end": 19.816666,
				  "confidence": 0.99907553
				},
				{
				  "word": "far",
				  "start": 19.816666,
				  "end": 20.053541,
				  "confidence": 0.52797174
				},
				{
				  "word": "long",
				  "start": 20.053541,
				  "end": 20.369373,
				  "confidence": 0.27643293
				},
				{
				  "word": "to",
				  "start": 20.369373,
				  "end": 20.52729,
				  "confidence": 0.9950035
				},
				{
				  "word": "realize",
				  "start": 20.52729,
				  "end": 21.02729,
				  "confidence": 0.9992204
				}
			  ]
			}
		  ]
		}
	  ]
	}
  }`

var MockSummarizeV1Response = `{
	"metadata": {
	  "transaction_key": "deprecated",
	  "request_id": "1bb89d91-53b8-499c-a352-2bccd75c8a7f",
	  "sha256": "5c3d7d6564176747187095038090a5bb44e703968c3836cfba29bf2c0c7d5dbe",
	  "created": "2023-07-03T19:21:59.006Z",
	  "duration": 21.687187,
	  "channels": 1,
	  "models": [
		"96a295ec-6336-43d5-b1cb-1e48b5e6d9a4"
	  ],
	  "model_info": {
		"96a295ec-6336-43d5-b1cb-1e48b5e6d9a4": {
		  "name": "general",
		  "version": "2023-02-22.3",
		  "arch": "base"
		}
	  }
	},
	"results": {
	  "channels": [
		{
		  "alternatives": [
			{
			  "transcript": "My illness has made very little difference to my scientific work. It was likely to have chosen theoretical physics. That is mainly your thought for which my physical disability is no handicap. I may be mentally disabled as well But if so, I am too far long to realize.",
			  "confidence": 0.99661994,
			  "words": [
				{
				  "word": "my",
				  "start": 0.439466,
				  "end": 0.59927183,
				  "confidence": 0.9942677,
				  "punctuated_word": "My"
				},
				{
				  "word": "illness",
				  "start": 0.59927183,
				  "end": 1.0992718,
				  "confidence": 0.99940467,
				  "punctuated_word": "illness"
				},
				{
				  "word": "has",
				  "start": 1.1585922,
				  "end": 1.3983009,
				  "confidence": 0.99580765,
				  "punctuated_word": "has"
				},
				{
				  "word": "made",
				  "start": 1.3983009,
				  "end": 1.7179126,
				  "confidence": 0.9998857,
				  "punctuated_word": "made"
				},
				{
				  "word": "very",
				  "start": 1.7179126,
				  "end": 2.0375242,
				  "confidence": 0.9997733,
				  "punctuated_word": "very"
				},
				{
				  "word": "little",
				  "start": 2.0375242,
				  "end": 2.4370387,
				  "confidence": 0.99988246,
				  "punctuated_word": "little"
				},
				{
				  "word": "difference",
				  "start": 2.4370387,
				  "end": 2.9164562,
				  "confidence": 0.98878103,
				  "punctuated_word": "difference"
				},
				{
				  "word": "to",
				  "start": 2.9164562,
				  "end": 3.076262,
				  "confidence": 0.91974664,
				  "punctuated_word": "to"
				},
				{
				  "word": "my",
				  "start": 3.076262,
				  "end": 3.3958735,
				  "confidence": 0.9985611,
				  "punctuated_word": "my"
				},
				{
				  "word": "scientific",
				  "start": 3.3958735,
				  "end": 3.8958735,
				  "confidence": 0.99998593,
				  "punctuated_word": "scientific"
				},
				{
				  "word": "work",
				  "start": 4.194903,
				  "end": 4.4346113,
				  "confidence": 0.99719137,
				  "punctuated_word": "work."
				},
				{
				  "word": "it",
				  "start": 4.914029,
				  "end": 5.1537375,
				  "confidence": 0.9884639,
				  "punctuated_word": "It"
				},
				{
				  "word": "was",
				  "start": 5.1537375,
				  "end": 5.3934464,
				  "confidence": 0.99965847,
				  "punctuated_word": "was"
				},
				{
				  "word": "likely",
				  "start": 5.3934464,
				  "end": 5.792961,
				  "confidence": 0.81015,
				  "punctuated_word": "likely"
				},
				{
				  "word": "to",
				  "start": 5.792961,
				  "end": 6.0326695,
				  "confidence": 0.9946583,
				  "punctuated_word": "to"
				},
				{
				  "word": "have",
				  "start": 6.0326695,
				  "end": 6.352281,
				  "confidence": 0.91991985,
				  "punctuated_word": "have"
				},
				{
				  "word": "chosen",
				  "start": 6.352281,
				  "end": 6.831699,
				  "confidence": 0.9999633,
				  "punctuated_word": "chosen"
				},
				{
				  "word": "theoretical",
				  "start": 6.831699,
				  "end": 7.331699,
				  "confidence": 0.99994814,
				  "punctuated_word": "theoretical"
				},
				{
				  "word": "physics",
				  "start": 7.630728,
				  "end": 8.110146,
				  "confidence": 0.99984443,
				  "punctuated_word": "physics."
				},
				{
				  "word": "that",
				  "start": 8.762314,
				  "end": 8.921488,
				  "confidence": 0.99965644,
				  "punctuated_word": "That"
				},
				{
				  "word": "is",
				  "start": 8.921488,
				  "end": 9.160248,
				  "confidence": 0.9083398,
				  "punctuated_word": "is"
				},
				{
				  "word": "mainly",
				  "start": 9.160248,
				  "end": 9.660248,
				  "confidence": 0.99661994,
				  "punctuated_word": "mainly"
				},
				{
				  "word": "your",
				  "start": 9.717356,
				  "end": 10.035703,
				  "confidence": 0.9974492,
				  "punctuated_word": "your"
				},
				{
				  "word": "thought",
				  "start": 10.035703,
				  "end": 10.535703,
				  "confidence": 0.96939987,
				  "punctuated_word": "thought"
				},
				{
				  "word": "for",
				  "start": 10.83157,
				  "end": 11.070331,
				  "confidence": 0.99951375,
				  "punctuated_word": "for"
				},
				{
				  "word": "which",
				  "start": 11.070331,
				  "end": 11.388678,
				  "confidence": 0.9998727,
				  "punctuated_word": "which"
				},
				{
				  "word": "my",
				  "start": 11.388678,
				  "end": 11.627438,
				  "confidence": 0.9997917,
				  "punctuated_word": "my"
				},
				{
				  "word": "physical",
				  "start": 11.627438,
				  "end": 12.127438,
				  "confidence": 0.9998171,
				  "punctuated_word": "physical"
				},
				{
				  "word": "disability",
				  "start": 12.5824795,
				  "end": 13.0824795,
				  "confidence": 0.9998565,
				  "punctuated_word": "disability"
				},
				{
				  "word": "is",
				  "start": 13.378347,
				  "end": 13.617107,
				  "confidence": 0.9026935,
				  "punctuated_word": "is"
				},
				{
				  "word": "no",
				  "start": 13.617107,
				  "end": 14.117107,
				  "confidence": 0.99946135,
				  "punctuated_word": "no"
				},
				{
				  "word": "handicap",
				  "start": 14.412975,
				  "end": 14.651735,
				  "confidence": 0.9574573,
				  "punctuated_word": "handicap."
				},
				{
				  "word": "i",
				  "start": 15.208843,
				  "end": 15.447603,
				  "confidence": 0.9988537,
				  "punctuated_word": "I"
				},
				{
				  "word": "may",
				  "start": 15.447603,
				  "end": 15.686363,
				  "confidence": 0.996563,
				  "punctuated_word": "may"
				},
				{
				  "word": "be",
				  "start": 15.686363,
				  "end": 15.925123,
				  "confidence": 0.996335,
				  "punctuated_word": "be"
				},
				{
				  "word": "mentally",
				  "start": 15.925123,
				  "end": 16.425123,
				  "confidence": 0.97879404,
				  "punctuated_word": "mentally"
				},
				{
				  "word": "disabled",
				  "start": 16.720991,
				  "end": 17.118925,
				  "confidence": 0.9996377,
				  "punctuated_word": "disabled"
				},
				{
				  "word": "as",
				  "start": 17.118925,
				  "end": 17.357685,
				  "confidence": 0.99978346,
				  "punctuated_word": "as"
				},
				{
				  "word": "well",
				  "start": 17.357685,
				  "end": 17.857685,
				  "confidence": 0.9997675,
				  "punctuated_word": "well"
				},
				{
				  "word": "but",
				  "start": 18.079582,
				  "end": 18.2375,
				  "confidence": 0.9989373,
				  "punctuated_word": "But"
				},
				{
				  "word": "if",
				  "start": 18.2375,
				  "end": 18.553333,
				  "confidence": 0.9977925,
				  "punctuated_word": "if"
				},
				{
				  "word": "so",
				  "start": 18.553333,
				  "end": 18.790207,
				  "confidence": 0.99553853,
				  "punctuated_word": "so,"
				},
				{
				  "word": "i",
				  "start": 19.185,
				  "end": 19.263958,
				  "confidence": 0.55681,
				  "punctuated_word": "I"
				},
				{
				  "word": "am",
				  "start": 19.263958,
				  "end": 19.500834,
				  "confidence": 0.8584063,
				  "punctuated_word": "am"
				},
				{
				  "word": "too",
				  "start": 19.500834,
				  "end": 19.816666,
				  "confidence": 0.99908173,
				  "punctuated_word": "too"
				},
				{
				  "word": "far",
				  "start": 19.816666,
				  "end": 20.053541,
				  "confidence": 0.53003824,
				  "punctuated_word": "far"
				},
				{
				  "word": "long",
				  "start": 20.053541,
				  "end": 20.369373,
				  "confidence": 0.2729079,
				  "punctuated_word": "long"
				},
				{
				  "word": "to",
				  "start": 20.369373,
				  "end": 20.52729,
				  "confidence": 0.99502146,
				  "punctuated_word": "to"
				},
				{
				  "word": "realize",
				  "start": 20.52729,
				  "end": 21.02729,
				  "confidence": 0.9992242,
				  "punctuated_word": "realize."
				}
			  ],
			  "summaries": [
				{
				  "summary": "My illness has made very little difference to my scientific work. It was likely to have chosen theoretical physics. I may be mentally disabled as well. But if so, I am too far long to realize.",
				  "start_word": 0,
				  "end_word": 49
				}
			  ]
			}
		  ]
		}
	  ]
	}
  }`

var MockSummarizeV2Response = `{
    "metadata":{
        "transaction_key":"deprecated",
        "request_id":"c856885a-d4b4-43de-8492-4a873e17cd22",
        "sha256":"5c3d7d6564176747187095038090a5bb44e703968c3836cfba29bf2c0c7d5dbe",
        "created":"2023-07-03T19:23:26.707Z",
        "duration":21.687187,
        "channels":1,
        "models":[
            "96a295ec-6336-43d5-b1cb-1e48b5e6d9a4"
        ],
        "model_info":{
            "96a295ec-6336-43d5-b1cb-1e48b5e6d9a4":{
                "name":"general",
                "version":"2023-02-22.3",
                "arch":"base"
            }
        }
    },
    "results":{
        "channels":[
            {
                "alternatives":[
                    {
                        "transcript":"My illness has made very little difference to my scientific work. It was likely to have chosen theoretical physics. That is mainly your thought for which my physical disability is no handicap. I may be mentally disabled as well But if so, I am too far long to realize.",
                        "confidence":0.9970703,
                        "words":[
                            {
                                "word":"my",
                                "start":0.439466,
                                "end":0.59927183,
                                "confidence":0.99560547,
                                "punctuated_word":"My"
                            },
                            {
                                "word":"illness",
                                "start":0.59927183,
                                "end":1.0992718,
                                "confidence":0.9995117,
                                "punctuated_word":"illness"
                            },
                            {
                                "word":"has",
                                "start":1.1585922,
                                "end":1.3983009,
                                "confidence":0.9970703,
                                "punctuated_word":"has"
                            },
                            {
                                "word":"made",
                                "start":1.3983009,
                                "end":1.7179126,
                                "confidence":1.0,
                                "punctuated_word":"made"
                            },
                            {
                                "word":"very",
                                "start":1.7179126,
                                "end":2.0375242,
                                "confidence":1.0,
                                "punctuated_word":"very"
                            },
                            {
                                "word":"little",
                                "start":2.0375242,
                                "end":2.4370387,
                                "confidence":1.0,
                                "punctuated_word":"little"
                            },
                            {
                                "word":"difference",
                                "start":2.4370387,
                                "end":2.9164562,
                                "confidence":0.99072266,
                                "punctuated_word":"difference"
                            },
                            {
                                "word":"to",
                                "start":2.9164562,
                                "end":3.076262,
                                "confidence":0.9277344,
                                "punctuated_word":"to"
                            },
                            {
                                "word":"my",
                                "start":3.076262,
                                "end":3.3958735,
                                "confidence":0.99853516,
                                "punctuated_word":"my"
                            },
                            {
                                "word":"scientific",
                                "start":3.3958735,
                                "end":3.8958735,
                                "confidence":1.0,
                                "punctuated_word":"scientific"
                            },
                            {
                                "word":"work",
                                "start":4.194903,
                                "end":4.4346113,
                                "confidence":0.9975586,
                                "punctuated_word":"work."
                            },
                            {
                                "word":"it",
                                "start":4.914029,
                                "end":5.1537375,
                                "confidence":0.9916992,
                                "punctuated_word":"It"
                            },
                            {
                                "word":"was",
                                "start":5.1537375,
                                "end":5.3934464,
                                "confidence":0.9995117,
                                "punctuated_word":"was"
                            },
                            {
                                "word":"likely",
                                "start":5.3934464,
                                "end":5.792961,
                                "confidence":0.8232422,
                                "punctuated_word":"likely"
                            },
                            {
                                "word":"to",
                                "start":5.792961,
                                "end":6.0326695,
                                "confidence":0.99560547,
                                "punctuated_word":"to"
                            },
                            {
                                "word":"have",
                                "start":6.0326695,
                                "end":6.352281,
                                "confidence":0.90478516,
                                "punctuated_word":"have"
                            },
                            {
                                "word":"chosen",
                                "start":6.352281,
                                "end":6.831699,
                                "confidence":1.0,
                                "punctuated_word":"chosen"
                            },
                            {
                                "word":"theoretical",
                                "start":6.831699,
                                "end":7.331699,
                                "confidence":1.0,
                                "punctuated_word":"theoretical"
                            },
                            {
                                "word":"physics",
                                "start":7.630728,
                                "end":8.030242,
                                "confidence":1.0,
                                "punctuated_word":"physics."
                            },
                            {
                                "word":"that",
                                "start":8.762314,
                                "end":8.921488,
                                "confidence":0.9995117,
                                "punctuated_word":"That"
                            },
                            {
                                "word":"is",
                                "start":8.921488,
                                "end":9.160248,
                                "confidence":0.9038086,
                                "punctuated_word":"is"
                            },
                            {
                                "word":"mainly",
                                "start":9.160248,
                                "end":9.660248,
                                "confidence":0.99658203,
                                "punctuated_word":"mainly"
                            },
                            {
                                "word":"your",
                                "start":9.717356,
                                "end":10.035703,
                                "confidence":0.9975586,
                                "punctuated_word":"your"
                            },
                            {
                                "word":"thought",
                                "start":10.035703,
                                "end":10.535703,
                                "confidence":0.96777344,
                                "punctuated_word":"thought"
                            },
                            {
                                "word":"for",
                                "start":10.83157,
                                "end":11.070331,
                                "confidence":0.9995117,
                                "punctuated_word":"for"
                            },
                            {
                                "word":"which",
                                "start":11.070331,
                                "end":11.388678,
                                "confidence":1.0,
                                "punctuated_word":"which"
                            },
                            {
                                "word":"my",
                                "start":11.388678,
                                "end":11.627438,
                                "confidence":1.0,
                                "punctuated_word":"my"
                            },
                            {
                                "word":"physical",
                                "start":11.627438,
                                "end":12.127438,
                                "confidence":1.0,
                                "punctuated_word":"physical"
                            },
                            {
                                "word":"disability",
                                "start":12.5824795,
                                "end":13.0824795,
                                "confidence":1.0,
                                "punctuated_word":"disability"
                            },
                            {
                                "word":"is",
                                "start":13.378347,
                                "end":13.617107,
                                "confidence":0.8979492,
                                "punctuated_word":"is"
                            },
                            {
                                "word":"no",
                                "start":13.617107,
                                "end":14.117107,
                                "confidence":0.9995117,
                                "punctuated_word":"no"
                            },
                            {
                                "word":"handicap",
                                "start":14.412975,
                                "end":14.651735,
                                "confidence":0.9584961,
                                "punctuated_word":"handicap."
                            },
                            {
                                "word":"i",
                                "start":15.208843,
                                "end":15.447603,
                                "confidence":0.99902344,
                                "punctuated_word":"I"
                            },
                            {
                                "word":"may",
                                "start":15.447603,
                                "end":15.686363,
                                "confidence":0.99658203,
                                "punctuated_word":"may"
                            },
                            {
                                "word":"be",
                                "start":15.686363,
                                "end":15.925123,
                                "confidence":0.99658203,
                                "punctuated_word":"be"
                            },
                            {
                                "word":"mentally",
                                "start":15.925123,
                                "end":16.425123,
                                "confidence":0.9785156,
                                "punctuated_word":"mentally"
                            },
                            {
                                "word":"disabled",
                                "start":16.720991,
                                "end":17.118925,
                                "confidence":0.9995117,
                                "punctuated_word":"disabled"
                            },
                            {
                                "word":"as",
                                "start":17.118925,
                                "end":17.357685,
                                "confidence":1.0,
                                "punctuated_word":"as"
                            },
                            {
                                "word":"well",
                                "start":17.357685,
                                "end":17.857685,
                                "confidence":1.0,
                                "punctuated_word":"well"
                            },
                            {
                                "word":"but",
                                "start":18.079582,
                                "end":18.2375,
                                "confidence":0.99902344,
                                "punctuated_word":"But"
                            },
                            {
                                "word":"if",
                                "start":18.2375,
                                "end":18.553333,
                                "confidence":0.9980469,
                                "punctuated_word":"if"
                            },
                            {
                                "word":"so",
                                "start":18.553333,
                                "end":18.790207,
                                "confidence":0.99560547,
                                "punctuated_word":"so,"
                            },
                            {
                                "word":"i",
                                "start":19.185,
                                "end":19.263958,
                                "confidence":0.55371094,
                                "punctuated_word":"I"
                            },
                            {
                                "word":"am",
                                "start":19.263958,
                                "end":19.500834,
                                "confidence":0.8544922,
                                "punctuated_word":"am"
                            },
                            {
                                "word":"too",
                                "start":19.500834,
                                "end":19.816666,
                                "confidence":0.99902344,
                                "punctuated_word":"too"
                            },
                            {
                                "word":"far",
                                "start":19.816666,
                                "end":20.053541,
                                "confidence":0.53759766,
                                "punctuated_word":"far"
                            },
                            {
                                "word":"long",
                                "start":20.053541,
                                "end":20.369373,
                                "confidence":0.26953125,
                                "punctuated_word":"long"
                            },
                            {
                                "word":"to",
                                "start":20.369373,
                                "end":20.52729,
                                "confidence":0.9951172,
                                "punctuated_word":"to"
                            },
                            {
                                "word":"realize",
                                "start":20.52729,
                                "end":21.02729,
                                "confidence":0.99902344,
                                "punctuated_word":"realize."
                            }
                        ]
                    }
                ]
            }
        ],
        "summary":{
            "short":"The speaker discusses their illness and how it has made little impact on their scientific work. They mention their physical disability and how it is likely related to theoretical physics. The speaker also acknowledges that they may be mentally disabled but it is too long before they realize their illness."
        }
    }
}`
