FROM rtuitlab/testmace:0.0.1-8

COPY ./Project /project
COPY ./waitWebApp.sh /waitWebApp.sh
RUN sed -i 's/\r$//' /waitWebApp.sh
RUN chmod +x /waitWebApp.sh

WORKDIR /out
CMD /waitWebApp.sh && testmace-cli /project --reporter=junit -e testmace