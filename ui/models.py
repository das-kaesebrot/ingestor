import secrets
from django.db import models


def generate_secret():
    return secrets.token_urlsafe(16)


class Project(models.Model):
    name = models.CharField(max_length=200)
    done_ingesting = models.BooleanField(default=False)


class User(models.Model):
    name = models.CharField(max_length=200)
    suffix = models.CharField(max_length=50)


class UserProjectLink(models.Model):
    user = models.ForeignKey(
        User, on_delete=models.CASCADE, related_name="project_links"
    )
    project = models.ForeignKey(
        Project, on_delete=models.CASCADE, related_name="user_links"
    )
    upload_secret = models.CharField(
        max_length=64, unique=True, default=generate_secret
    )

    class Meta:
        unique_together = ("user", "project")


class Device(models.Model):
    user_project_link = models.ForeignKey(
        UserProjectLink, on_delete=models.CASCADE, related_name="devices"
    )
    name = models.CharField(max_length=100)
    suffix = models.CharField(max_length=50)

    class Meta:
        unique_together = ("user_project_link", "name")


class File(models.Model):
    device = models.ForeignKey(Device, on_delete=models.CASCADE, related_name="files")
    file = models.FileField(upload_to="uploads/")
    filename = models.CharField(max_length=255)
    mime_type = models.CharField(max_length=100, null=True, blank=True)
    size = models.BigIntegerField(null=True, blank=True)
    uploaded_at = models.DateTimeField(auto_now_add=True)
    kept = models.BooleanField(default=True)

    def save(self, *args, **kwargs):
        if not self.filename:
            self.filename = self.file.name
        if not self.size:
            self.size = self.file.size
        super().save(*args, **kwargs)
