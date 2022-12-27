import connexion
import six
from typing import Dict
from typing import Tuple
from typing import Union

from openapi_server.models.nwdaf_ml_model_prov_subsc import NwdafMLModelProvSubsc  # noqa: E501
from openapi_server.models.problem_details import ProblemDetails  # noqa: E501
from openapi_server import util


def create_nwdafml_model_provision_subcription(nwdaf_ml_model_prov_subsc=None):  # noqa: E501
    """Create a new Individual NWDAF ML Model Provision Subscription resource.

     # noqa: E501

    :param nwdaf_ml_model_prov_subsc: 
    :type nwdaf_ml_model_prov_subsc: dict | bytes

    :rtype: Union[NwdafMLModelProvSubsc, Tuple[NwdafMLModelProvSubsc, int], Tuple[NwdafMLModelProvSubsc, int, Dict[str, str]]
    """
    if connexion.request.is_json:
        nwdaf_ml_model_prov_subsc = NwdafMLModelProvSubsc.from_dict(connexion.request.get_json())  # noqa: E501
    return nwdaf_ml_model_prov_subsc
