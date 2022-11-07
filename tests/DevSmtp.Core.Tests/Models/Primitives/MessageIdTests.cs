using DevSmtp.Core.Models;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace DevSmtp.Core.Tests.Models.Primitives
{
    [TestClass]
    public class MessageIdTests
    {
        [TestMethod]
        public void Ctor_WhenValueIsValid_ItShouldCreateUserId()
        {
            // Arrange
            var value = "someuserid";

            // Act
            var results = new MessageId(value);

            // Assert
            Assert.IsNotNull(results);
            Assert.AreEqual(value, results.Value);
        }

        [TestMethod]
        public void From_WhenValueIsValid_ItShouldCreateUserId()
        {
            // Arrange
            var value = "someuserid";

            // Act
            var results = MessageId.From(value);

            // Assert
            Assert.IsNotNull(results);
            Assert.AreEqual(value, results.Value);
        }

        [TestMethod]
        public void Ctor_WhenValueIsNotValid_ItShouldThrowFormatException()
        {
            // Arrange
            var value = "";

            // Act
            try
            {
                _ = new MessageId(value);
                Assert.Fail("FormatException expected");
            }
            catch (FormatException)
            {
                // expected...
            }
        }
    }
}
