{{ define "reg" }}

<!DOCTYPE html>
<html>
    <head>
        <title>Регистрация</title>
        <meta charset="utf-8">
        <link rel="stylesheet" href="/web/css/style.css">
        <script type="text/javascript" src="/web/js/reg_verification.js"></script>
        <meta http-equiv="Cache-control" content="no-cache">
        <meta http-equiv="Pragma" content="no-cache">
        <meta http-equiv="Expires" content="-1">
        <style>
            form br {
                margin-top: 5px;
                margin-bottom: 5px;
            }
        </style>
    </head>
    <body>
        <div class="main_div">
            {{ template "header" }}
            <form action="/reg_verification/" method="POST" id="reg_form">
                <input type="text" name="username" placeholder="Имя пользователя" required> <br>
                <input type="password" name="password1" id="password1" placeholder="Пароль"  required> <br>
                <input type="password" name="password2" id="password2" placeholder="Повтор пароль" required> <br>
                <err_tag id="err_tag" style="color: red;"></err_tag> <br>
                <input type="submit" value="Подтвердить" onclick="regVerification(event);">
            </form>
        </div>
    </body>
</html>

{{ end }}