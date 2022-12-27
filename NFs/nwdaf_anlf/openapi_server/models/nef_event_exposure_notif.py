# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.nef_event_notification import NefEventNotification
from openapi_server import util

from openapi_server.models.nef_event_notification import NefEventNotification  # noqa: E501

class NefEventExposureNotif(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, notif_id=None, event_notifs=None):  # noqa: E501
        """NefEventExposureNotif - a model defined in OpenAPI

        :param notif_id: The notif_id of this NefEventExposureNotif.  # noqa: E501
        :type notif_id: str
        :param event_notifs: The event_notifs of this NefEventExposureNotif.  # noqa: E501
        :type event_notifs: List[NefEventNotification]
        """
        self.openapi_types = {
            'notif_id': str,
            'event_notifs': List[NefEventNotification]
        }

        self.attribute_map = {
            'notif_id': 'notifId',
            'event_notifs': 'eventNotifs'
        }

        self.notif_id = notif_id
        self.event_notifs = event_notifs

    @classmethod
    def from_dict(cls, dikt) -> 'NefEventExposureNotif':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The NefEventExposureNotif of this NefEventExposureNotif.  # noqa: E501
        :rtype: NefEventExposureNotif
        """
        return util.deserialize_model(dikt, cls)

    @property
    def notif_id(self):
        """Gets the notif_id of this NefEventExposureNotif.


        :return: The notif_id of this NefEventExposureNotif.
        :rtype: str
        """
        return self._notif_id

    @notif_id.setter
    def notif_id(self, notif_id):
        """Sets the notif_id of this NefEventExposureNotif.


        :param notif_id: The notif_id of this NefEventExposureNotif.
        :type notif_id: str
        """
        if notif_id is None:
            raise ValueError("Invalid value for `notif_id`, must not be `None`")  # noqa: E501

        self._notif_id = notif_id

    @property
    def event_notifs(self):
        """Gets the event_notifs of this NefEventExposureNotif.


        :return: The event_notifs of this NefEventExposureNotif.
        :rtype: List[NefEventNotification]
        """
        return self._event_notifs

    @event_notifs.setter
    def event_notifs(self, event_notifs):
        """Sets the event_notifs of this NefEventExposureNotif.


        :param event_notifs: The event_notifs of this NefEventExposureNotif.
        :type event_notifs: List[NefEventNotification]
        """
        if event_notifs is None:
            raise ValueError("Invalid value for `event_notifs`, must not be `None`")  # noqa: E501
        if event_notifs is not None and len(event_notifs) < 1:
            raise ValueError("Invalid value for `event_notifs`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._event_notifs = event_notifs
