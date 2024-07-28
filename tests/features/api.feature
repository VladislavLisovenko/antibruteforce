# file: features/api.feature

# http://anti-bruteforce:8080/

Feature: Check user
    Scenario: Check if auth data available
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=test&password=pass&ip=127.0.0.1"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":false}
        """

    Scenario: Reset auth
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=testReset&password=passReset&ip=127.0.0.2"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":false}
        """
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=testReset&password=passReset&ip=127.0.0.2"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":false}
        """
        When I send "GET" request to "http://anti-bruteforce:8080/reset?ip=127.0.0.2&login=testReset"
        Then The response code should be 200
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=testReset&password=passReset&ip=127.0.0.2"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":true}
        """

    Scenario: Method not allowed
        When I send "GET" request to "http://anti-bruteforce:8080/blackList?ip=127.0.0.1"
        Then The response code should be 405

    Scenario: Check ip in blacklist
        When I send "POST" request to "http://anti-bruteforce:8080/blackList?ip=127.0.0.3/32"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":true}
        """
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=testBl&password=passBl&ip=127.0.0.3"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":false}
        """
        When I send "DELETE" request to "http://anti-bruteforce:8080/blackList?ip=127.0.0.3/32"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":true}
        """
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=testBl&password=passBl&ip=127.0.0.3"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":false}
        """

    Scenario: Check ip in whitelist
        When I send "POST" request to "http://anti-bruteforce:8080/whiteList?ip=127.0.0.4/32"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":true}
        """
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=testWl&password=passBl&ip=127.0.0.4"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":true}
        """
        When I send "GET" request to "http://anti-bruteforce:8080/check?login=testWl&password=passBl&ip=127.0.0.4"
        Then The response code should be 200
        And The response should match text:
        """
        {"ok":true}
        """