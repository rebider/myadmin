@echo off
rd /s/q E:\WEB\0226\myadmin\views
md E:\WEB\0226\myadmin\views
xcopy /s /e /y  dist  E:\WEB\0226\myadmin\views\
pause
