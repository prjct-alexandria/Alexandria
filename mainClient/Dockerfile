# source: https://mherman.org/blog/dockerizing-a-react-app/

# pull official base image
FROM node:16.13.1-alpine

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install app dependencies
COPY package.json ./
COPY package-lock.json ./
RUN npm install
RUN npm install react-scripts@3.4.1 -g

EXPOSE 3000

# add app
COPY . ./

# start app
CMD ["npm", "start"]