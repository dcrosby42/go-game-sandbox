#version 330
precision highp float;

uniform mat4 MVP_MATRIX;
uniform mat4 MV_MATRIX;
in vec3 VERTEX_POSITION;
in vec3 VERTEX_NORMAL;
in vec2 VERTEX_UV_0;

out vec3 vs_position;
out vec3 vs_eye_normal;
out vec2 vs_uv_0;

void main()
{
  vs_position = VERTEX_POSITION;
  vs_uv_0 = VERTEX_UV_0;
  vs_eye_normal = normalize(mat3(MV_MATRIX) * VERTEX_NORMAL);
  gl_Position = MVP_MATRIX * vec4(VERTEX_POSITION, 1.0);
}
