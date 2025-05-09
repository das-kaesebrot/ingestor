from django.shortcuts import render
from django.http import HttpResponse, HttpRequest
from django.views.decorators.http import require_GET
from django_htmx.middleware import HtmxDetails

class HtmxHttpRequest(HttpRequest):
    htmx: HtmxDetails

@require_GET
def index(request: HtmxHttpRequest) -> HttpResponse:
    return render(request, "index.html")
