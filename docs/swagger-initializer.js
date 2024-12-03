window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    urls: [
      // skip addition
      {url: "proto/addition/v1/service.swagger.json", name: "addition"},
      {url: "proto/subtraction/v1/service.swagger.json", name: "subtraction"},
      {url: "proto/multiplication/v1/service.swagger.json", name: "multiplication"},
      {url: "proto/division/v1/service.swagger.json", name: "division"},
    ],
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
