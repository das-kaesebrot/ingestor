from django.urls import path

from . import views

urlpatterns = [
    path("", views.index, name="index"),
    path("favicon.ico", views.favicon),
    path("projects", views.project_list),
    path(
        "upload/<str:upload_secret>/",
        views.upload_file,
        name="upload_file",
    ),
    path("project/<int:project_id>/files/", views.project_files, name="project_files"),
    # Project CRUD
    path("projects/", views.project_list, name="project_list"),
    path("projects/add/", views.project_add, name="project_add"),
    path("projects/<int:project_id>/", views.project_detail, name="project_detail"),
    path("projects/<int:project_id>/edit/", views.project_edit, name="project_edit"),
    path(
        "projects/<int:project_id>/delete/", views.project_delete, name="project_delete"
    ),
    # User management
    path(
        "projects/<int:project_id>/add_user/",
        views.project_add_user,
        name="project_add_user",
    ),
    path("user_links/<int:link_id>/edit/", views.user_edit, name="user_edit"),
    path("user_links/<int:link_id>/delete/", views.user_delete, name="user_delete"),
]
