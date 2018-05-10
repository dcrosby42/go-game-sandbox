#version 330

precision highp float;
uniform sampler2D MATERIAL_TEX_0;

uniform mat4 MV_MATRIX;
uniform mat4 MVP_MATRIX;
uniform vec4 MATERIAL_DIFFUSE;

in vec3 vs_position;
in vec3 vs_eye_normal;
in vec2 vs_uv_0;

out vec4 frag_color;

void main()
{
  vec4 LIGHT_POSITION = vec4(-100.0, 100.0, 100.0, 1.0);
  vec4 LIGHT_DIFFUSE = vec4(1.0);

  vec4 eye_vertex = MV_MATRIX * vec4(vs_position, 1.0);
  vec3 s = normalize(vec3(LIGHT_POSITION-eye_vertex));
  float diffuseIntensity = max(dot(s,vs_eye_normal), 0.0);

  vec4 ambient = vec4(0.1,0.1,0.1,0.1);


  frag_color = MATERIAL_DIFFUSE * texture(MATERIAL_TEX_0, vs_uv_0) * diffuseIntensity;
  frag_color = frag_color + ambient;
}
