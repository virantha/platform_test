import requests
import json
import pytest
import random

class TestRoutes(object):
    """Run py.test on this"""

    def setup(self):
        self.HOST = "http://localhost:8080/message/route"
            
    def _send(self):
        r = requests.post(self.HOST, data=json.dumps(self.payload))
        return r
            

    def _assert_error(self, r, errnum, raised_message=None):
        with pytest.raises(requests.HTTPError):
            r.raise_for_status()
        assert r.status_code == errnum
        if raised_message:
            assert raised_message in json.loads(r.text)["error"]

    def test_error_malformed(self):
        self.payload = {"message":"test", "recipients":"ldkj"}
        r = self._send()
        self._assert_error(r, 400, "Malformed request")

    def test_error_nomessage(self):
        self.payload = {"message":"", 'recipient':"1231231234"}
        r = self._send()
        self._assert_error(r,400, "Message cannot be empty")

    def test_error_norecipients(self):
        self.payload = {"message":"rest", 'recipient':"[]"}
        r = self._send()
        self._assert_error(r,400, "Recipients list cannot be empty")


    def test_error_toomany(self):
        self.payload = {"message":"rest", 'recipients':["1231231234"]*5000}
        r = self._send()
        r.raise_for_status()

        self.payload = {"message":"rest", 'recipients':["1231231234"]*5001}
        r = self._send()
        self._assert_error(r,400, "but maximum allowed is 5000")

    def test_error_invalidphone(self):

        self.payload = {"message":"rest", 'recipients':["1231231234","123456789"]}
        r = self._send()
        self._assert_error(r,400, "Invalid phone")

    def _check_return(self, result, expected):
        """result and expected should be dicts
        """
        assert expected["message"] == result["message"]

        # Make a hash of the result, indexed by ip
        result_ip = {}
        for route in result["routes"]:
            ip = route["ip"]
            recipients = route["recipients"]
            if ip in result_ip:
                result_ip[ip].extend(recipients)
            else:
                result_ip[ip] = recipients

        for route in expected["routes"]:
            ip = route["ip"]
            for recipient in route["recipients"]:
                assert recipient in result_ip[ip]

    def _gen_phone_numbers(self, n):
        """ n is how many to generate"""
        number_list = []
        for i in range(n):
            l = [str(random.randint(0,9)) for x in range(10)]
            number_list.append(''.join(l))
        return number_list

    def test_post_simple(self):
        self.payload = { "message": "Testing 1",
                         "recipients": ["1231231234", "2342342345"]
                       }
        expected = {"message": "Testing 1",
                    "routes": [
                            {   "ip": "10.0.1.1",
                                "recipients": ["1231231234"],
                            },
                            {   "ip": "10.0.1.2",
                                "recipients": ["2342342345"],
                            },
                            ]
                        }

        r = self._send()
        r.raise_for_status()
        result = json.loads(r.text)
        self._check_return(result, expected)

    def test_post_11(self):
        recipients = self._gen_phone_numbers(11)
        self.payload = { "message": "Testing 2",
                         "recipients": recipients
                       }

        expected = {"message": "Testing 2",
                    "routes": [
                            {   "ip": "10.0.3.1",
                                "recipients": recipients[:10]
                            },
                            {   "ip": "10.0.1.1",
                                "recipients": [recipients[10]]
                            },
                            ]
                        }

        print expected
        r = self._send()
        r.raise_for_status()
        result = json.loads(r.text)
        self._check_return(result, expected)

    def test_post_17(self):
        recipients = self._gen_phone_numbers(17)
        self.payload = { "message": "Testing 3",
                         "recipients": recipients
                       }

        expected = {"message": "Testing 3",
                    "routes": [
                            {   "ip": "10.0.3.1",
                                "recipients": recipients[:10]
                            },
                            {   "ip": "10.0.2.1",
                                "recipients": recipients[10:15]
                            },
                            {   "ip": "10.0.1.1",
                                "recipients": [recipients[15]]
                            },
                            {   "ip": "10.0.1.2",
                                "recipients": [recipients[16]]
                            },
                            ]
                        }

        print expected
        r = self._send()
        r.raise_for_status()
        result = json.loads(r.text)
        self._check_return(result, expected)

    def test_post_51(self):
        recipients = self._gen_phone_numbers(51)
        self.payload = { "message": "Testing 4",
                         "recipients": recipients
                       }

        expected = {"message": "Testing 4",
                    "routes": [
                            {   "ip": "10.0.4.1",
                                "recipients": recipients[:25]
                            },
                            {   "ip": "10.0.4.2",
                                "recipients": recipients[25:50]
                            },
                            {   "ip": "10.0.1.1",
                                "recipients": [recipients[50]]
                            },
                            ]
                        }

        print expected
        r = self._send()
        r.raise_for_status()
        result = json.loads(r.text)
        self._check_return(result, expected)
