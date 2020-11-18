@ECHO OFF
IF %1==NO-CACHE docker build --no-cache --tag atlas-wrg:latest .
IF NOT %1==NO-CACHE docker build --tag atlas-wrg:latest .