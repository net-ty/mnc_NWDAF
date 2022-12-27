# coding: utf-8

from __future__ import absolute_import
import unittest

from flask import json
from six import BytesIO

from openapi_server.models.analytics_data import AnalyticsData  # noqa: E501
from openapi_server.models.any_of_event_id_any_ofstring import AnyOfEventIdAnyOfstring  # noqa: E501
from openapi_server.models.event_filter import EventFilter  # noqa: E501
from openapi_server.models.event_reporting_requirement import EventReportingRequirement  # noqa: E501
from openapi_server.models.problem_details import ProblemDetails  # noqa: E501
from openapi_server.models.problem_details_analytics_info_request import ProblemDetailsAnalyticsInfoRequest  # noqa: E501
from openapi_server.models.target_ue_information import TargetUeInformation  # noqa: E501
from openapi_server.test import BaseTestCase


class TestNWDAFAnalyticsDocumentController(BaseTestCase):
    """NWDAFAnalyticsDocumentController integration test stubs"""

    def test_get_nwdaf_analytics(self):
        """Test case for get_nwdaf_analytics

        Read a NWDAF Analytics
        """
        query_string = [('event-id', openapi_server.EventId()),
                        ('ana-req', {'key': openapi_server.EventReportingRequirement()}),
                        ('event-filter', {'key': openapi_server.EventFilter()}),
                        ('supported-features', 'supported_features_example'),
                        ('tgt-ue', {'key': openapi_server.TargetUeInformation()})]
        headers = { 
            'Accept': 'application/json',
            'Authorization': 'Bearer special-key',
        }
        response = self.client.open(
            '/nnwdaf-analyticsinfo/v1/analytics',
            method='GET',
            headers=headers,
            query_string=query_string)
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))


if __name__ == '__main__':
    unittest.main()
