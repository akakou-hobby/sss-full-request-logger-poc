from flask import Flask, request, render_template
app = Flask(__name__)


@app.route("/")
def index():
    name = request.args.get('name', 'Taro')
    return render_template('index.html', name=name)

if __name__ == "__main__":
    app.run(host='0.0.0.0')
