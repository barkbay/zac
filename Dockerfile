FROM centos:7
#from golang:1.7.5-alpine

COPY zac /zac

CMD /zac server
