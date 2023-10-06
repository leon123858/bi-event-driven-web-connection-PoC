FROM node:alpine
# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
# Install app dependencies
COPY package.json /usr/src/app/
RUN npm install
# Bundle app source
COPY . /usr/src/app
RUN npm run build
# Execute
EXPOSE 3000
ENV HOSTNAME 0.0.0.0
ENV PORT 3000
# entrypoint
CMD [ "npm", "start" ]