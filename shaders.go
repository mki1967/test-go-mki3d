package main

var vertexShader = `
#version 330
layout (location = 0) in vec3 position;
void main() {
    gl_Position = vec4(position, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
out vec4 outputColor;
void main() {
    outputColor = vec4(1.0, 0.0, 1.0, 1.0);
}
` + "\x00"
