import connexion
import six
from typing import Dict
from typing import Tuple
from typing import Union

from openapi_server.models.context_data import ContextData  # noqa: E501
from openapi_server.models.context_id_list import ContextIdList  # noqa: E501
from openapi_server.models.problem_details import ProblemDetails  # noqa: E501
from openapi_server.models.requested_context import RequestedContext  # noqa: E501
from openapi_server import util


def get_nwdaf_context(context_ids, req_context=None):  # noqa: E501
    """Get context information related to analytics subscriptions.

     # noqa: E501

    :param context_ids: Identifies specific context information related to analytics subscriptions.
    :type context_ids: dict | bytes
    :param req_context: Identfies the type(s) of the analytics context information the consumer wishes to receive. 
    :type req_context: dict | bytes

    :rtype: Union[ContextData, Tuple[ContextData, int], Tuple[ContextData, int, Dict[str, str]]
    """
    if connexion.request.is_json:
        context_ids =  ContextIdList.from_dict(connexion.request.get_json())  # noqa: E501
    if connexion.request.is_json:
        req_context =  RequestedContext.from_dict(connexion.request.get_json())  # noqa: E501
    return 'do some magic!'
