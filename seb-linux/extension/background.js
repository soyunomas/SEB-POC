var SEB_CONFIG_KEY = "2ac94554b0aefff6b6abc5de67a767c7e66711baa919e217e53889478f5e77ea";
var SEB_BEK = "1e9b2524a1b021966a337ab2881a6c42ddf510be13f56cf80bba3fc9fcb476eb";
var SEB_USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 SafeExamBrowser/3.3.0";

chrome.webRequest.onBeforeSendHeaders.addListener(
  function(details) {
    var url = details.url.split("#")[0];

    var configKeyHash = SHA256(url + SEB_CONFIG_KEY);
    var requestHash = SHA256(url + SEB_BEK);

    var headers = details.requestHeaders;
    headers.push({name: "X-SafeExamBrowser-ConfigKeyHash", value: configKeyHash});
    headers.push({name: "X-SafeExamBrowser-RequestHash", value: requestHash});

    for (var i = 0; i < headers.length; i++) {
      if (headers[i].name.toLowerCase() === "user-agent") {
        headers[i].value = SEB_USER_AGENT;
      }
    }

    console.log("[SEB] URL:", url);
    console.log("[SEB] ConfigKeyHash:", configKeyHash);
    console.log("[SEB] RequestHash:", requestHash);

    return {requestHeaders: headers};
  },
  {urls: ["<all_urls>"]},
  ["blocking", "requestHeaders"]
);

console.log("[SEB-Linux] Extension loaded. ConfigKey:", SEB_CONFIG_KEY);
