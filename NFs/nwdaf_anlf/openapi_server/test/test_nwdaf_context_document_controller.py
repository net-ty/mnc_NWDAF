# coding: utf-8

from __future__ import absolute_import
import unittest

from flask import json
from six import BytesIO

from openapi_server.models.context_data import ContextData  # noqa: E501
from openapi_server.models.context_id_list import ContextIdList  # noqa: E501
from openapi_server.models.problem_details import ProblemDetails  # noqa: E501
from openapi_server.models.requested_context import RequestedContext  # noqa: E501
from openapi_server.test import BaseTestCase


class TestNWDAFContextDocumentController(BaseTestCase):
    """NWDAFContextDocumentController integration test stubs"""

    def test_get_nwdaf_context(self):
        """Test case for get_nwdaf_context

        Get context information related to analytics subscriptions.
        """
        query_string = [('context-ids', {'key': openapi_server.ContextIdList()}),
                        ('req-context', {'key': openapi_server.RequestedContext()})]
        headers = { 
            'Accept': 'application/json',
            'Authorization': 'Bearer special-key',
        }
        response = self.client.open(
            '/nnwdaf-analyticsinfo/v1/context',
            method='GET',
            headers=headers,
            query_string=query_string)
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))


if __name__ == '__main__':
    unittest.main()
