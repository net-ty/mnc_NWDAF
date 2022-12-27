# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.nwdaf_event import NwdafEvent
from openapi_server.models.nwdaf_failure_code import NwdafFailureCode
from openapi_server import util

from openapi_server.models.nwdaf_event import NwdafEvent  # noqa: E501
from openapi_server.models.nwdaf_failure_code import NwdafFailureCode  # noqa: E501

class FailureEventInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, event=None, failure_code=None):  # noqa: E501
        """FailureEventInfo - a model defined in OpenAPI

        :param event: The event of this FailureEventInfo.  # noqa: E501
        :type event: NwdafEvent
        :param failure_code: The failure_code of this FailureEventInfo.  # noqa: E501
        :type failure_code: NwdafFailureCode
        """
        self.openapi_types = {
            'event': NwdafEvent,
            'failure_code': NwdafFailureCode
        }

        self.attribute_map = {
            'event': 'event',
            'failure_code': 'failureCode'
        }

        self.event = event
        self.failure_code = failure_code

    @classmethod
    def from_dict(cls, dikt) -> 'FailureEventInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The FailureEventInfo of this FailureEventInfo.  # noqa: E501
        :rtype: FailureEventInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def event(self):
        """Gets the event of this FailureEventInfo.


        :return: The event of this FailureEventInfo.
        :rtype: NwdafEvent
        """
        return self._event

    @event.setter
    def event(self, event):
        """Sets the event of this FailureEventInfo.


        :param event: The event of this FailureEventInfo.
        :type event: NwdafEvent
        """
        if event is None:
            raise ValueError("Invalid value for `event`, must not be `None`")  # noqa: E501

        self._event = event

    @property
    def failure_code(self):
        """Gets the failure_code of this FailureEventInfo.


        :return: The failure_code of this FailureEventInfo.
        :rtype: NwdafFailureCode
        """
        return self._failure_code

    @failure_code.setter
    def failure_code(self, failure_code):
        """Sets the failure_code of this FailureEventInfo.


        :param failure_code: The failure_code of this FailureEventInfo.
        :type failure_code: NwdafFailureCode
        """
        if failure_code is None:
            raise ValueError("Invalid value for `failure_code`, must not be `None`")  # noqa: E501

        self._failure_code = failure_code
