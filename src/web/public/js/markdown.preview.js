$('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
  if(e.target.hash === '#preview') {
    var mdTest = document.getElementById('md-text').value;
    document.getElementById('md-preview').innerHTML = marked(mdTest, { sanitize: true });
  }
});
