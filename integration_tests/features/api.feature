Feature: Test event service API
	Check API methods

	Scenario: API Create and Get Event
		Given there is user "test_user"
		And there is server "localhost:8080"
		When I create event
		"""
		{
			"title":"test_event_create_get",
			"text":"Test event: API Create check",
			"startTime":"2019-10-01T00:00:00Z",
			"endTime":"2019-10-02T00:00:00Z"
		}
		"""
		And I get event by previous id
		Then Events should be the same

	Scenario: API Update and Get Event
		Given there is user "test_user"
		And there is server "localhost:8080"
		When I create event
		"""
		{
			"title":"test_event_create_update",
			"text":"Test event: API Create Update check",
			"startTime":"2019-10-03T00:00:00Z",
			"endTime":"2019-10-04T00:00:00Z"
		}
		"""
		And I update created event
		"""
		{
			"title":"test_event_update",
			"text":"Test event: API Create check",
			"startTime":"2019-10-05T00:00:00Z",
			"endTime":"2019-10-06T00:00:00Z"
		}
		"""
		And I get event by previous id
		Then Events should be the same
