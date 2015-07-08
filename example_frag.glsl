#version 330


    #define FOG_DISABLED -1
    #define FOG_LINEAR 0
    #define FOG_EXP 1
    #define FOG_EXP2 2

    struct Fog
    {
        vec4 color;
        float start;
        float end;
        float density;
        int mode;
    };

    float getFogFactor(Fog fog, float distance)
    {
        float factor = 0.0;

        if (fog.mode == FOG_DISABLED) {
            return 0;
        } else if (fog.mode == FOG_LINEAR) {
            factor = (fog.end - distance) / (fog.end - fog.start);
        } else if (fog.mode == FOG_EXP) {
            factor = exp(-fog.density * distance);
        } else if (fog.mode == FOG_EXP2) {
            factor = exp(-pow(fog.density * distance, 2.0));
        }

        return 1.0 - clamp(factor, 0.0, 1.0);
    }

    vec4 getFog(Fog fog, float distance)
    {
        if (fog.mode == FOG_DISABLED) {
            return vec4(1);
        }

        float factor = getFogFactor(fog, distance);

        factor = 1.0 - clamp(factor, 0.0, 1.0);

        return fog.color * vec4(factor);
    }

    out fragColor0 out;
uniform Fog fog;

    in vec4 position;

    vec4 get0(color)
    {
        float fogCoord = abs(position.z/position.w);
        return mix(color, fog.color, getFogFactor(fog, fogCoord));
    }


    struct BaseLight
    {
        vec3 color;
        float intensity;
    };

    struct Attenuation
    {
        float constant;
        float linear;
        float exponent;
    };

    vec4 CalcLight(BaseLight base, vec3 direction, vec3 normal, vec3 worldPos, vec3 eyePos, float specularIntensity, float specularPower)
    {
    	float diffuseFactor = dot(normal, -direction);

    	vec4 diffuseColor = vec4(0,0,0,0);
    	vec4 specularColor = vec4(0,0,0,0);

    	if (diffuseFactor > 0) {
    		diffuseColor = vec4(base.color, 1.0) * base.intensity * diffuseFactor;

    		vec3 eyeDir = normalize(eyePos - worldPos);

    		// Phong Lighting Model
    		vec3 reflectionDir = normalize(reflect(direction, normal));
    		float specularFactor = dot(eyeDir, reflectionDir);

    		// Almost Phong, but cheaper! ;)
    		// vec3 halfDir = normalize(eyeDir - direction);
    		// float specularFactor = dot(halfDir, normal);

    		specularFactor = pow(specularFactor, specularPower);

    		if (specularFactor > 0) {
    			specularColor = vec4(base.color, 1.0) * specularIntensity * specularFactor;
    		}
    	}

    	return diffuseColor + specularColor;
    }

    vec4 CalcPointLight(PointLight pointLight, vec3 normal, vec3 worldPos, vec3 eyePos, float specularIntensity, float specularPower)
    {
    	vec3 lightDir = worldPos - pointLight.position;
    	float distance2Point = length(lightDir);

    	if (distance2Point > pointLight.range) {
    	    return vec4(0, 0, 0, 0);
    	}

    	lightDir = normalize(lightDir);

    	vec4 color =  CalcLight(pointLight.base, lightDir, normal, worldPos, eyePos, specularIntensity, specularPower);

    	float attenuation = pointLight.attenuation.constant +
    						pointLight.attenuation.linear * distance2Point +
    						pointLight.attenuation.exponent * distance2Point * distance2Point +
    						0.0001f; // Make calc division by 0 safe.

    	return color / attenuation;
    }

    vec4 CalcDirectionalLight(DirectionalLight directionalLight, vec3 normal, vec3 worldPos, vec3 eyePos, float specularIntensity, float specularPower)
    {
    	return CalcLight(directionalLight.base, directionalLight.direction, normal, worldPos, eyePos, specularIntensity, specularPower);
    }

    vec4 CalcSpotLight(SpotLight spotLight, vec3 normal, vec3 worldPos, vec3 eyePos, float specularIntensity, float specularPower)
    {
    	vec3 dir = normalize(worldPos - spotLight.pointLight.position);
    	float factor = dot(dir, spotLight.direction);

    	vec4 color = vec4(0, 0, 0, 0);

    	if (factor > spotLight.cutoff) {
    		color = CalcPointLight(spotLight.pointLight, normal, worldPos, eyePos, specularIntensity, specularPower);
    		// Fuzzy edges!
    		color = color  * (1.0 - (1.0 - factor) / (1.0 - spotLight.cutoff));
    	}

    	return color;
    }


in vec2 texCoord;

uniform vec4 color;
uniform sampler2D diffuse;

void main()
{
    vec4 out = color * texture(diffuse, texCoord.xy);

    out = get0(out);out = get1(out);

    fragColor0 = out;
}


