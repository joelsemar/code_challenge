import requests
import unittest
from datetime import datetime
import random
import json
import time

"""
This test module is more of an integration test suite, but using the python unit test lib
for code organization

Assumes ua.code.challenge binary available 127.0.0.1 port 8080
"""

API_LOCATION = "http://127.0.0.1:8080%s"


class UACodeChallengeTest(unittest.TestCase):

    def test_can_create_message(self):
        test_user = "joel"
        test_text = self.random_text()

        resp = self.create_message(test_text, test_user)
        self.assertEqual(resp.status_code, 201)
        self.assertTrue(self.message_exists_in_chat(test_text, test_user))

        # after one read, the message should no longer show up
        self.assertFalse(self.message_exists_in_chat(test_text, test_user))

    def test_message_expires_after_timeout(self):
        test_user = "joel"
        test_text = self.random_text()

        self.create_message(test_text, test_user, timeout=1)
        # little wiggle room
        time.sleep(2)
        self.assertFalse(self.message_exists_in_chat(test_text, test_user))

    def test_can_hit_10_tps(self):
        # basically set 5 messages for five users, then read each of their chats,
        # this makes 10 requests (5 writes, 5 reads)
        # do the loop 100 times, should be 1000 requests

        start = datetime.utcnow()
        for i in range(100):
            names = ["larry", "curly", "moe", "joel", "paul"]
            for name in names:
                self.create_message(self.random_text(), name)
                requests.get(self.url("/chat/" + name))
        seconds = (datetime.utcnow() - start).total_seconds()
        self.assertLess(seconds, 100)
        print "Executed %s requests in %s seconds" % (1000, seconds)

    def create_message(self, text, user, timeout=60):
        payload = {
            "text": text,
            "username": user,
            "timeout": timeout
        }
        return requests.post(self.url("/chat"), json=payload)

    def url(self, path):
        return API_LOCATION % path

    def random_text(self, length=32):
        alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
        return ''.join([random.choice(alphanum) for i in range(length)])

    def message_exists_in_chat(self, message, username):
        messages = requests.get(self.url("/chat/" + username)).content
        messages = json.loads(messages)

        return message in [m["text"] for m in messages]


if __name__ == '__main__':
    unittest.main()
