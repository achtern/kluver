#version 330

layout (location = 0) in vec3 inPosition;
layout (location = 1) in vec2 inTexCoord;
layout (location = 2) in vec3 inNormal;


out vec2 texCoord;
out vec3 normal;
out vec3 worldPos;
out vec4 shadowMapCoord;
out vec4 position;
uniform mat4 model;
uniform mat4 modelView;
uniform mat4 MVP;
uniform mat4 shadowMatrix;

void main ()
{
  gl_Position = MVP * vec4(inPosition, 1.0);

  vec2 texCoord  = inTexCoord;
  vec3 normal  = (model * vec4(inNormal, 0.0)).xyz;
  vec3 worldPos  = (model * vec4(inPosition, 1.0)).xyz;
  vec4 shadowMapCoord  = shadowMatrix * vec4(inPosition, 1.0);

  
    vec4 position  = modelView * vec4(position, 1.0);

}
