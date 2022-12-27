# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.pdu_session_status import PduSessionStatus
from openapi_server import util

from openapi_server.models.pdu_session_status import PduSessionStatus  # noqa: E501

class PduSessionInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, n4_sess_id=None, sess_inactive_timer=None, pdu_sess_status=None):  # noqa: E501
        """PduSessionInfo - a model defined in OpenAPI

        :param n4_sess_id: The n4_sess_id of this PduSessionInfo.  # noqa: E501
        :type n4_sess_id: str
        :param sess_inactive_timer: The sess_inactive_timer of this PduSessionInfo.  # noqa: E501
        :type sess_inactive_timer: int
        :param pdu_sess_status: The pdu_sess_status of this PduSessionInfo.  # noqa: E501
        :type pdu_sess_status: PduSessionStatus
        """
        self.openapi_types = {
            'n4_sess_id': str,
            'sess_inactive_timer': int,
            'pdu_sess_status': PduSessionStatus
        }

        self.attribute_map = {
            'n4_sess_id': 'n4SessId',
            'sess_inactive_timer': 'sessInactiveTimer',
            'pdu_sess_status': 'pduSessStatus'
        }

        self.n4_sess_id = n4_sess_id
        self.sess_inactive_timer = sess_inactive_timer
        self.pdu_sess_status = pdu_sess_status

    @classmethod
    def from_dict(cls, dikt) -> 'PduSessionInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PduSessionInfo of this PduSessionInfo.  # noqa: E501
        :rtype: PduSessionInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def n4_sess_id(self):
        """Gets the n4_sess_id of this PduSessionInfo.


        :return: The n4_sess_id of this PduSessionInfo.
        :rtype: str
        """
        return self._n4_sess_id

    @n4_sess_id.setter
    def n4_sess_id(self, n4_sess_id):
        """Sets the n4_sess_id of this PduSessionInfo.


        :param n4_sess_id: The n4_sess_id of this PduSessionInfo.
        :type n4_sess_id: str
        """

        self._n4_sess_id = n4_sess_id

    @property
    def sess_inactive_timer(self):
        """Gets the sess_inactive_timer of this PduSessionInfo.

        indicating a time in seconds.  # noqa: E501

        :return: The sess_inactive_timer of this PduSessionInfo.
        :rtype: int
        """
        return self._sess_inactive_timer

    @sess_inactive_timer.setter
    def sess_inactive_timer(self, sess_inactive_timer):
        """Sets the sess_inactive_timer of this PduSessionInfo.

        indicating a time in seconds.  # noqa: E501

        :param sess_inactive_timer: The sess_inactive_timer of this PduSessionInfo.
        :type sess_inactive_timer: int
        """

        self._sess_inactive_timer = sess_inactive_timer

    @property
    def pdu_sess_status(self):
        """Gets the pdu_sess_status of this PduSessionInfo.


        :return: The pdu_sess_status of this PduSessionInfo.
        :rtype: PduSessionStatus
        """
        return self._pdu_sess_status

    @pdu_sess_status.setter
    def pdu_sess_status(self, pdu_sess_status):
        """Sets the pdu_sess_status of this PduSessionInfo.


        :param pdu_sess_status: The pdu_sess_status of this PduSessionInfo.
        :type pdu_sess_status: PduSessionStatus
        """

        self._pdu_sess_status = pdu_sess_status
