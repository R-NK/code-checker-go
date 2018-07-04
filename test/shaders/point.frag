#version 330 core
out vec4 fragment;

in vec3 vertexColor;
in vec2 TexCord;

uniform sampler2D texture1;
uniform sampler2D texture2;

void main()
{
	fragment = mix(texture(texture1, TexCord), texture(texture2, TexCord), 0.2);
}