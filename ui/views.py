from django.shortcuts import render
from django.http import HttpResponse, HttpRequest
from django.views.decorators.http import require_GET
from django_htmx.middleware import HtmxDetails
from django.core.paginator import Paginator
from .models import Person, Project

BASE_TEMPLATE = "_base.html"
PARTIAL_TEMPLATE = "_partial.html"

class HtmxHttpRequest(HttpRequest):
    htmx: HtmxDetails

# determine whether the request is a full http request or just a partial htmx one, render content accordingly
def get_base_template_name(request: HtmxHttpRequest) -> str:
    if request.htmx:
        return PARTIAL_TEMPLATE
    return BASE_TEMPLATE


@require_GET
def index(request: HtmxHttpRequest) -> HttpResponse:
    return render(
        request,
        "index.html",
        {
            "base_template": get_base_template_name(request),
        },
    )

@require_GET
def favicon(request: HtmxHttpRequest) -> HttpResponse:
    return HttpResponse(
        (
            '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">'
            + '<text y=".9em" font-size="90">🛂</text>'
            + "</svg>"
        ),
        content_type="image/svg+xml",
    )


@require_GET
def project_list(request: HtmxHttpRequest) -> HttpResponse:    
    # Standard Django pagination
    page_num = request.GET.get("page", "1")
    page = Paginator(object_list=Project.objects.all(), per_page=10).get_page(page_num)        

    return render(
        request,
        "project_list.html",
        {
            "base_template": get_base_template_name(request),
            "page": page,
        },
    )
    