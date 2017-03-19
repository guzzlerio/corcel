
filename="${1}"

#Ginkgo has a custom style to start the tests where as GoConvey simply requires a pure go Test method
sed -i -E 's/var _ = Describe\("([^")]+)".*$/func Test\1\(t *testing.T\){\n\tConvey("\1", t, func(){/g' ${filename}
echo "}" >> ${filename}

#The different test structure methods of Ginkgo replaced with the only one for GoConvey
sed -i -E 's/Describe\(\"|Context\(\"|It\(\"/Convey(\"/g' ${filename}

#Assertions with GoConvey start with So(...), Ginkgo has an Expect method aswell as others
# but this project only uses the Expect method.
sed -i -E 's/Expect/So/g' ${filename}

#Migrating those assertions for To(Equal(...))
sed -i -E 's/\)\.To\(Equal\(([^\)]+)\)\)/, ShouldEqual, \1\)/g' ${filename}

#Migrating those assertions for To(BeNil())
sed -i -E 's/\)\.To\(BeNil\(\)\)/, ShouldBeNil)/g' ${filename}

#This feels like a useful thing to remember for the future.  Multiline grep
#grep -Poz "(?s)import \([^\)]+\)" ${filename}
