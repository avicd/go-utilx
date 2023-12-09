package xmlx

import (
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	xml := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<!DOCTYPE mapper PUBLIC \"-//mybatis.org//DTD Mapper 3.0//EN\"\n" +
		"        \"http://mybatis.org/dtd/mybatis-3-mapper.dtd\">\n" +
		"<mapper namespace=\"MapperBuilder\">\n" +
		"    <sql id=\"sometable\">\n" +
		"        ${prefix}Table\n" +
		"    </sql>\n" +
		"\n" +
		"    <sql id=\"someinclude-t2\">\n" +
		"        from\n" +
		"        <include refid=\"${include_target}\"/>\n" +
		"    </sql>\n" +
		"\n" +
		"    <select id=\"select\" resultType=\"map\">\n" +
		"        select\n" +
		"        field1, field2, field3\n" +
		"\n" +
		"        <sql id=\"someinclude\">\n" +
		"            from\n" +
		"            <include refid=\"${include_target}\"/>\n" +
		"        </sql>\n" +
		"        <include refid=\"someinclude\">\n" +
		"            <property name=\"prefix\" value=\"Some\"/>\n" +
		"            <property name=\"include_target\" value=\"sometable\"/>\n" +
		"        </include>\n" +
		"        <sql id=\"someinclude-s2\">\n" +
		"            from\n" +
		"            <include refid=\"${include_target}\"/>\n" +
		"        </sql>\n" +
		"    </select>\n" +
		"    <sql id=\"someinclude-t3\">\n" +
		"        from\n" +
		"        <include refid=\"${include_target}\"/>\n" +
		"    </sql>\n" +
		"</mapper>"
	xml1 := "<sml><sss/><sdfsdfds/></sml>"
	node, _ := Parse(strings.NewReader(xml))
	node, _ = Parse(strings.NewReader(xml1))
	s := node.FindOne("//sss")
	s.Remove()
	println(s)
}
