@import fog.slib;

#---VERTEX---#
layout (location = 0) in vec3 inPosition;
layout (location = 1) in vec2 inTexCoord;
layout (location = 2) in vec3 inNormal;


@require mat4 model;
@require mat4 modelView;
@require mat4 MVP;
@require mat4 shadowMatrix;

void main ()
{
  gl_Position = MVP * vec4(inPosition, 1.0);

  @provide vec2 texCoord = inTexCoord;
  @provide vec3 normal = (model * vec4(inNormal, 0.0)).xyz;
  @provide vec3 worldPos = (model * vec4(inPosition, 1.0)).xyz;
  @provide vec4 shadowMapCoord = shadowMatrix * vec4(inPosition, 1.0);

  @yield;
}
#---END---#
#---FRAGMENT---#
@request vec2 texCoord;

@require vec4 color;
@require sampler2D diffuse;

void main()
{
    vec4 out = color * texture(diffuse, texCoord.xy);

    @yield out;

    @write(0) out;
}


#---END---#