# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.context_type import ContextType
from openapi_server import util

from openapi_server.models.context_type import ContextType  # noqa: E501

class RequestedContext(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, contexts=None):  # noqa: E501
        """RequestedContext - a model defined in OpenAPI

        :param contexts: The contexts of this RequestedContext.  # noqa: E501
        :type contexts: List[ContextType]
        """
        self.openapi_types = {
            'contexts': List[ContextType]
        }

        self.attribute_map = {
            'contexts': 'contexts'
        }

        self.contexts = contexts

    @classmethod
    def from_dict(cls, dikt) -> 'RequestedContext':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The RequestedContext of this RequestedContext.  # noqa: E501
        :rtype: RequestedContext
        """
        return util.deserialize_model(dikt, cls)

    @property
    def contexts(self):
        """Gets the contexts of this RequestedContext.

        List of analytics context types.  # noqa: E501

        :return: The contexts of this RequestedContext.
        :rtype: List[ContextType]
        """
        return self._contexts

    @contexts.setter
    def contexts(self, contexts):
        """Sets the contexts of this RequestedContext.

        List of analytics context types.  # noqa: E501

        :param contexts: The contexts of this RequestedContext.
        :type contexts: List[ContextType]
        """
        if contexts is None:
            raise ValueError("Invalid value for `contexts`, must not be `None`")  # noqa: E501
        if contexts is not None and len(contexts) < 1:
            raise ValueError("Invalid value for `contexts`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._contexts = contexts
