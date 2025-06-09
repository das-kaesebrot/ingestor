from django.db import models

class Person(models.Model):
    name = models.CharField(max_length=128)
    
class Project(models.Model):
    name = models.CharField(max_length=30)
    created_at = models.DateField()
    members = models.ManyToManyField(Person)
