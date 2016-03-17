<html>
<head>
    <title></title>
    <link rel="stylesheet" href="./bootstrap.min.css" />
</head>
<body>
<form action="/login2" method="POST">
    <input type="checkbox" name="interest" value="football">Football
    <input type="checkbox" name="interest" value="basketball">Basketball
    <input type="checkbox" name="interest" value="tennis">Tennis
    Username2:<input type="text" name="username">
    Password2:<input type="password" name="password">
    <input type="hidden" name="token" value="{{.}}">
    <input type="submit" value="Login">
</form>
</body>
</html>