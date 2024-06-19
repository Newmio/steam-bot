chrome.webRequest.onAuthRequired.addListener(
    function(details, callback) {
      callback({
        authCredentials: {username: "yggdjocl", password: "vajq3n53awr1"}
      });
    },
    {urls: ["<all_urls>"]},
    ['blocking']
  );
  