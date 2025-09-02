from django.shortcuts import get_object_or_404, redirect, render
from django.http import HttpResponse, HttpRequest
from django.views.decorators.http import require_GET
from django_htmx.middleware import HtmxDetails
from django.core.paginator import Paginator

from .forms import FileUploadForm, UserForm, ProjectForm
from .models import Device, File, Project, UserProjectLink

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
def project_list(request):
    projects = Project.objects.all().order_by("id")
    return render(request, "project/list.html", {"projects": projects})


def upload_file(request, upload_secret, device_id):
    uplink = get_object_or_404(UserProjectLink, upload_secret=upload_secret)
    device = get_object_or_404(Device, id=device_id, user_project_link=uplink)

    if request.method == "POST":
        form = FileUploadForm(request.POST, request.FILES)
        if form.is_valid():
            file_instance = File(
                device=device,
                file=form.cleaned_data["file"],
                mime_type=form.cleaned_data["file"].content_type,
            )
            file_instance.save()
            return redirect("project/files.html", project_id=uplink.project.id)
    else:
        form = FileUploadForm()

    return render(request, "upload.html", {"form": form, "device": device})


def project_files(request, project_id):
    files = File.objects.filter(
        device__user_project_link__project_id=project_id
    ).order_by("-uploaded_at")
    return render(request, "project/files.html", {"files": files})


def project_add(request):
    if request.method == "POST":
        form = ProjectForm(request.POST)
        if form.is_valid():
            form.save()
            return redirect("project_list")
    else:
        form = ProjectForm()
    return render(request, "project/add.html", {"form": form})


def project_detail(request, project_id):
    project = get_object_or_404(Project, id=project_id)
    user_links = project.user_links.select_related("user").all()
    return render(
        request,
        "project/detail.html",
        {
            "project": project,
            "user_links": user_links,
        },
    )


def project_add_user(request, project_id):
    project = get_object_or_404(Project, id=project_id)
    if request.method == "POST":
        form = UserForm(request.POST)
        if form.is_valid():
            user = form.save()
            UserProjectLink.objects.create(
                user=user, project=project
            )  # secret auto-generated
            return redirect("project_detail", project_id=project.id)
    else:
        form = UserForm()
    return render(request, "project/add_user.html", {"form": form, "project": project})


# -----------------
# Project CRUD
# -----------------
def project_edit(request, project_id):
    project = get_object_or_404(Project, id=project_id)
    if request.method == "POST":
        form = ProjectForm(request.POST, instance=project)
        if form.is_valid():
            form.save()
            return redirect("project_detail", project_id=project.id)
    else:
        form = ProjectForm(instance=project)
    return render(request, "project/edit.html", {"form": form, "project": project})


def project_delete(request, project_id):
    project = get_object_or_404(Project, id=project_id)
    if request.method == "POST":
        project.delete()
        return redirect("project_list")
    return render(request, "project/delete.html", {"project": project})


# -----------------
# User management
# -----------------
def user_edit(request, link_id):
    link = get_object_or_404(UserProjectLink, id=link_id)
    user = link.user
    if request.method == "POST":
        form = UserForm(request.POST, instance=user)
        if form.is_valid():
            form.save()
            return redirect("project_detail", project_id=link.project.id)
    else:
        form = UserForm(instance=user)
    return render(request, "user/edit.html", {"form": form, "user": user, "link": link})


def user_delete(request, link_id):
    link = get_object_or_404(UserProjectLink, id=link_id)
    project_id = link.project.id
    if request.method == "POST":
        link.delete()
        return redirect("project_detail", project_id=project_id)
    return render(request, "user/delete.html", {"link": link})
