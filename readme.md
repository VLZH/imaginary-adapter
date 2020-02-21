# imaginary

Adapter for [h2non/imaginary](https://github.com/h2non/imaginary)

I don't like original "HTTP API semantic" of imaginary. I have need link to image like this: `http://somehost.com/uploads/image.png?method=fit&height=300&widht=300`.

In case when I do not define method I would be to get an original image by Nginx or Caddy.

In usually I use `imaginary` like service in Docker Swarm, all traffic from clients going to Nginx that proxying requests to `imaginary` if `method` is defined. But I was forced to use ugly rewrite rules in `Nginx` configuration, with this adapter I can just proxy requests to `imaginary-adapter` service.

# Configuration

For configure adapter you must to use environment variables:

- **ADAPTER_PORT** - adapter server port (default: 9000)
- **ADAPTER_HOST** - adapter server host (default: 0.0.0.0)
- **ADAPTER_IMAGINARY_HOST** - url to imaginary host
- **ADAPTER_FILE_PATH_PREFIX** - this part will cutting from url
- **ADAPTER_DEFAULT_TYPE** - the default type of response image. If it defined as `auto`, imaginary will look to the header `Accept` and returns an image with a specific format. If `ADAPTER_DEFAULT_TYPE` is not defined type will be "" and imaginary will return image with original format

## How to run in Docker Compose:

```yml
services:
  imaginary:
    image: h2non/imaginary:latest
    volumes:
      - uploads:/mnt/data
    environment:
      PORT: 9000
    command: -enable-url-source -mount /mnt/data
    networks:
      - overlay
  imaginary-adapter:
    image: vlzhvlzh/imaginary-adapter
    environment:
      ADAPTER_IMAGINARY_HOST: http://imaginary:9000
      ADAPTER_FILE_PATH_PREFIX: /uploads
      ADAPTER_PORT: 9000
    ports:
      - "9000:9000"
    networks:
      - overlay
```

With the configuration above you can to open image by URL: `http://locahost:9000/uploads/someimage.png?method=crop&height=300&width=300`

# Disclaimer

Images in `./uploads`:

- ./uploads/1.jpeg - [unsplash](https://unsplash.com/photos/s9xm11AEpqU)
- ./uploads/2.jpeg - [unsplash](https://unsplash.com/photos/wfUBXRu1uSU)
- ./uploads/3.jpeg - [unsplash](https://unsplash.com/photos/aKK8q8Vl23U)

# TODO`s:

- Parameters validation
- Support of all methods
