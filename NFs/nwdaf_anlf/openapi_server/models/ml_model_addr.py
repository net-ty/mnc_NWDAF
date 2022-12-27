# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class MLModelAddr(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, m_l_model_url=None, ml_file_fqdn=None):  # noqa: E501
        """MLModelAddr - a model defined in OpenAPI

        :param m_l_model_url: The m_l_model_url of this MLModelAddr.  # noqa: E501
        :type m_l_model_url: str
        :param ml_file_fqdn: The ml_file_fqdn of this MLModelAddr.  # noqa: E501
        :type ml_file_fqdn: str
        """
        self.openapi_types = {
            'm_l_model_url': str,
            'ml_file_fqdn': str
        }

        self.attribute_map = {
            'm_l_model_url': 'mLModelUrl',
            'ml_file_fqdn': 'mlFileFqdn'
        }

        self.m_l_model_url = m_l_model_url
        self.ml_file_fqdn = ml_file_fqdn

    @classmethod
    def from_dict(cls, dikt) -> 'MLModelAddr':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The MLModelAddr of this MLModelAddr.  # noqa: E501
        :rtype: MLModelAddr
        """
        return util.deserialize_model(dikt, cls)

    @property
    def m_l_model_url(self):
        """Gets the m_l_model_url of this MLModelAddr.

        String providing an URI formatted according to RFC 3986.  # noqa: E501

        :return: The m_l_model_url of this MLModelAddr.
        :rtype: str
        """
        return self._m_l_model_url

    @m_l_model_url.setter
    def m_l_model_url(self, m_l_model_url):
        """Sets the m_l_model_url of this MLModelAddr.

        String providing an URI formatted according to RFC 3986.  # noqa: E501

        :param m_l_model_url: The m_l_model_url of this MLModelAddr.
        :type m_l_model_url: str
        """

        self._m_l_model_url = m_l_model_url

    @property
    def ml_file_fqdn(self):
        """Gets the ml_file_fqdn of this MLModelAddr.

        The FQDN of the ML Model file.  # noqa: E501

        :return: The ml_file_fqdn of this MLModelAddr.
        :rtype: str
        """
        return self._ml_file_fqdn

    @ml_file_fqdn.setter
    def ml_file_fqdn(self, ml_file_fqdn):
        """Sets the ml_file_fqdn of this MLModelAddr.

        The FQDN of the ML Model file.  # noqa: E501

        :param ml_file_fqdn: The ml_file_fqdn of this MLModelAddr.
        :type ml_file_fqdn: str
        """

        self._ml_file_fqdn = ml_file_fqdn
