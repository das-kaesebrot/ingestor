from django import forms

from .models import Project, User


class FileUploadForm(forms.Form):
    file = forms.FileField()
    device = forms.IntegerField()


class ProjectForm(forms.ModelForm):
    class Meta:
        model = Project
        fields = ["name", "done_ingesting"]


class UserForm(forms.ModelForm):
    class Meta:
        model = User
        fields = ["name", "suffix"]
