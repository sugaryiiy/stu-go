FROM 121.36.61.64:8888/common/alpine:3.20.0
LABEL authors="xiaogu"
ADD stu-go /
ENTRYPOINT ["/stu-go"]