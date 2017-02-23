FROM alpine:3.3
ADD grpctestapp /
CMD ["/grpctestapp", "counter"]