{{ define "footer" }}
    <!-- jQuery -->
    <script src="{{.URLPathPrefix}}/lib/startbootstrap-sb-admin-2/vendor/jquery/jquery.min.js"></script> 

    <!-- Bootstrap Core JavaScript 
    <script src="{{.URLPathPrefix}}/lib/startbootstrap-sb-admin-2/vendor/bootstrap/js/bootstrap.min.js"></script> -->
    <!-- Our modified Bootstrap -->
    <script src="{{.URLPathPrefix}}/lib/bootstrap/js/bootstrap.min.js"></script>
    
    <!-- Bootstrap-Dialog -->
    <script src="https://cdnjs.cloudflare.com/ajax/libs/bootstrap3-dialog/1.34.7/js/bootstrap-dialog.min.js"></script>
    
    <!-- Metis Menu Plugin JavaScript -->
    <script src="{{.URLPathPrefix}}/lib/startbootstrap-sb-admin-2/vendor/metisMenu/metisMenu.min.js"></script>

    <!-- Custom Theme JavaScript -->
    <script src="{{.URLPathPrefix}}/lib/startbootstrap-sb-admin-2/dist/js/sb-admin-2.js"></script>

</body>

</html>
{{ end }}